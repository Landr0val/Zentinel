import Hashids from 'hashids';

const SECRET = process.env.NEXT_PUBLIC_HASHID_SECRET || 'zentinel_default_secret_key_2024';
const MIN_LENGTH = 12;
const ALPHABET = 'abcdefghijklmnopqrstuvwxyz0123456789';

const hashids = new Hashids(SECRET, MIN_LENGTH, ALPHABET);

export function createSlug(databaseId: string): string {
  try {
    const hex = databaseId.replace(/-/g, '');
    return hashids.encodeHex(hex);
  } catch (error) {
    console.error('Error creating slug');
    return databaseId;
  }
}

export function getIdFromSlug(slug: string): string | null {
  try {
    const decodedHex = hashids.decodeHex(slug);
    return decodedHex || null;
  } catch (error) {
    return null;
  }
}
