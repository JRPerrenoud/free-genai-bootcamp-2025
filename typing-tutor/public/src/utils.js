// Utility functions for the Spanish typing tutor

/**
 * Normalizes Spanish text by removing accents and converting to lowercase
 * This allows users to type with an English keyboard
 * 
 * @param {string} text - The text to normalize
 * @returns {string} - Normalized text without accents
 */
function normalizeSpanish(text) {
    if (!text) return '';
    return text.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLowerCase();
}
