/**
 * Convert base64 audio string to blob URL
 * @param base64 - Base64 encoded audio string (with or without data URI prefix)
 * @returns Object URL for the audio blob
 */
export function base64ToBlobUrl(base64: string): string {
  let mime = "audio/wav";
  let pureBase64 = base64;

  if (base64.startsWith("data:")) {
    const split = base64.split(",");
    const header = split[0];
    pureBase64 = split[1];
    mime = header.match(/data:(.*?);base64/)?.[1] || mime;
  }

  const byteCharacters = atob(pureBase64);
  const byteNumbers = new Array(byteCharacters.length);

  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i);
  }

  const byteArray = new Uint8Array(byteNumbers);
  const blob = new Blob([byteArray], { type: mime });

  return URL.createObjectURL(blob);
}

