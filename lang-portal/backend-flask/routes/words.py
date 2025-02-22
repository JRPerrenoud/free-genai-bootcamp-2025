from flask import request, jsonify, g
from flask_cors import cross_origin
import json

def load(app):
  # Endpoint: GET /api/words with pagination (50 words per page)
  @app.route('/api/words', methods=['GET'])
  @cross_origin()
  def get_words():
    try:
      cursor = app.db.cursor()

      # Get the current page number from query parameters (default is 1)
      page = int(request.args.get('page', 1))
      # Ensure page number is positive
      page = max(1, page)
      words_per_page = 50
      offset = (page - 1) * words_per_page

      # Get sorting parameters from the query string
      sort_by = request.args.get('sort_by', 'english')  # Default to sorting by 'english'
      order = request.args.get('order', 'asc')  # Default to ascending order
      group = request.args.get('group')  # Get group filter

      # Validate sort_by and order
      valid_columns = ['english', 'spanish', 'correct_count', 'wrong_count']
      if sort_by not in valid_columns:
        sort_by = 'english'
      if order not in ['asc', 'desc']:
        order = 'asc'

      # Base query for words
      base_query = '''
        SELECT DISTINCT w.id, w.english, w.spanish, 
            COALESCE(r.correct_count, 0) AS correct_count,
            COALESCE(r.wrong_count, 0) AS wrong_count
        FROM words w
        LEFT JOIN word_reviews r ON w.id = r.word_id
      '''

      # Base query for count
      count_query = '''
        SELECT COUNT(DISTINCT w.id)
        FROM words w
      '''

      # Add group filter if specified
      params = []
      if group:
        group_join = '''
          JOIN word_groups wg ON w.id = wg.word_id
          JOIN groups g ON wg.group_id = g.id
          WHERE g.name = ?
        '''
        base_query += group_join
        count_query += group_join
        params.append(group)

      # Add sorting and pagination to the base query
      base_query += f' ORDER BY {sort_by} {order} LIMIT ? OFFSET ?'
      params_with_limit = params.copy()
      params_with_limit.extend([words_per_page, offset])

      # Get total words count
      cursor.execute(count_query, params)
      total_words = cursor.fetchone()[0]
      total_pages = (total_words + words_per_page - 1) // words_per_page

      # Get words for current page
      cursor.execute(base_query, params_with_limit)

      words = cursor.fetchall()

      # Format the response
      words_data = []
      for word in words:
        words_data.append({
          "id": word["id"],
          "english": word["english"],
          "spanish": word["spanish"],
          "correct_count": word["correct_count"],
          "wrong_count": word["wrong_count"]
        })

      return jsonify({
        "words": words_data,
        "total_pages": total_pages,
        "current_page": page,
        "total_words": total_words
      })

    except Exception as e:
      return jsonify({"error": str(e)}), 500
    finally:
      app.db.close()

  # Endpoint: GET /api/words/:id to get a single word with its details
  @app.route('/api/words/<int:word_id>', methods=['GET'])
  @cross_origin()
  def get_word(word_id):
    try:
      cursor = app.db.cursor()
      
      # Query to fetch the word and its details
      cursor.execute('''
        SELECT w.id, w.english, w.spanish,
               COALESCE(r.correct_count, 0) AS correct_count,
               COALESCE(r.wrong_count, 0) AS wrong_count,
               GROUP_CONCAT(DISTINCT g.id || '::' || g.name) as groups
        FROM words w
        LEFT JOIN word_reviews r ON w.id = r.word_id
        LEFT JOIN word_groups wg ON w.id = wg.word_id
        LEFT JOIN groups g ON wg.group_id = g.id
        WHERE w.id = ?
        GROUP BY w.id
      ''', (word_id,))
      
      word = cursor.fetchone()
      
      if not word:
        return jsonify({"error": "Word not found"}), 404
      
      # Parse the groups string into a list of group objects
      groups = []
      if word["groups"]:
        group_strings = word["groups"].split(',')
        for group_string in group_strings:
          group_id, group_name = group_string.split('::')
          groups.append({
            "id": int(group_id),
            "name": group_name
          })
      
      return jsonify({
        "id": word["id"],
        "english": word["english"],
        "spanish": word["spanish"],
        "correct_count": word["correct_count"],
        "wrong_count": word["wrong_count"],
        "groups": groups
      })
      
    except Exception as e:
      return jsonify({"error": str(e)}), 500
    finally:
      app.db.close()