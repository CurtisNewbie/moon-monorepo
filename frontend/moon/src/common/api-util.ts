const TOKEN = "token";

let emptyTokenCallback = null;
const AuthenticationCookieName = "Gatekeeper_Authorization";

export function onEmptyToken(callback) {
  emptyTokenCallback = callback;
}

export function setToken(token: string) {
  if (token === null) {
    localStorage.removeItem(TOKEN);
    document.cookie = "";
  } else {
    localStorage.setItem(TOKEN, token);
    if (document.cookie.includes(AuthenticationCookieName)) {
      setEventPumpCookie(token);
    }
  }
}

export function getToken() {
  let tkn = localStorage.getItem(TOKEN);
  if (!tkn && emptyTokenCallback) {
    emptyTokenCallback();
  }
  return tkn;
}

export function setEventPumpCookie(token = null) {
  if (token === null) {
    token = getToken();
  }
  if (token) {
    document.cookie =
      AuthenticationCookieName + " = Bearer " + token + "; path = /event-pump;";
  }
}
