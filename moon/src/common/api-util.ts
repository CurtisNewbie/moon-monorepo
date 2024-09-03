const TOKEN = "token";

let emptyTokenCallback = null;

export function onEmptyToken(callback) {
  emptyTokenCallback = callback;
}

export function setToken(token: string) {
  if (token === null) localStorage.removeItem(TOKEN);
  else {
    localStorage.setItem(TOKEN, token);
  }
}

export function getToken() {
  let tkn = localStorage.getItem(TOKEN);
  if (!tkn && emptyTokenCallback) {
    emptyTokenCallback();
  }
  return tkn;
}
