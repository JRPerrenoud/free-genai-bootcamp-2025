## Setting up the database

```sh
invoke init-db
```

This will do the following:
- Create the words.db (SQLite database)
- Create database tables using SQL definitions from `sql/setup/` directory
- Import seed data from the `seed/` directory (including word lists and study activities)
- NOTE: This is a desctructive action and will delete the existing database if it exists. Only do this if you want to reset the database.

## Clearing the database

Simply delete the `words.db` to clear the entire database.

## Running the backend API

```sh
python app.py 
```

This will start the Flask app on port `5000`

## Project Structure

- `app.py` - Main Flask application entry point
- `lib/db.py` - Database connection and initialization logic
- `routes/` - API endpoint definitions
- `seed/` - JSON files containing initial data
- `sql/setup/` - SQL files for table creation
- `tasks.py` - Invoke tasks for database initialization
