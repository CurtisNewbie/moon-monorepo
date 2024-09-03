
export function copyToClipboard(t: string) {
    if (!t) return;

    // src: https://stackoverflow.com/questions/400212/how-do-i-copy-to-the-clipboard-in-javascript
    var textarea = document.createElement("textarea");
    textarea.textContent = t;
    textarea.style.position = "fixed"; // Prevent scrolling to bottom of page in Microsoft Edge.
    document.body.appendChild(textarea);
    textarea.select();

    try {
        return document.execCommand("copy"); // Security exception may be thrown by some browsers.
    } catch (ex) {
        console.warn("Copy to clipboard failed.", ex);
    } finally {
        document.body.removeChild(textarea);
    }
}