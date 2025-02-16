/**
 * Gets the thumbnail URL for a study activity
 * First tries the provided thumbnail_url if it's an HTTP URL
 * Then tries to use a matching PNG based on the activity name
 * Falls back to placeholder.png if no match is found
 */
export function getActivityThumbnailUrl(activityName: string, providedUrl?: string): string {
  // If we have a valid HTTP URL, use it
  if (providedUrl?.startsWith('http')) {
    return providedUrl
  }

  // Convert activity name to lowercase and replace spaces with hyphens
  const normalizedName = activityName.toLowerCase().replace(/\s+/g, '-')
  
  // Return the normalized name as the image path
  return `/${normalizedName}.png`
}
