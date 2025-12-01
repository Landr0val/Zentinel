import Hashids from 'hashids';

const SECRET = process.env.NEXT_PUBLIC_HASHID_SECRET || 'zentinel_default_secret_key_2024';
const MIN_LENGTH = 12;
const ALPHABET = 'abcdefghijklmnopqrstuvwxyz0123456789';

const hashids = new Hashids(SECRET, MIN_LENGTH, ALPHABET);

export function createSlug(databaseId: string): string {
  try {
    const numericId = parseInt(databaseId, 10);
    if (!isNaN(numericId)) {
      return hashids.encode(numericId);
    }
    const hashValue = Array.from(databaseId).reduce((acc, char) => acc + char.charCodeAt(0), 0);
    return hashids.encode(hashValue);
  } catch (error) {
    console.error('Error creating slug:', error);
    return databaseId;
  }
}

export function getIdFromSlug(slug: string): string | null {
  try {
    const decoded = hashids.decode(slug);
    if (decoded.length > 0) {
      return decoded[0].toString();
    }
    return null;
  } catch (error) {
    console.error('Error decoding slug:', error);
    return null;
  }
}

export function removeSlug(_slug: string): void {
}

export function clearAllSlugs(): void {
}