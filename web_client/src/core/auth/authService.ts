import axios, { AxiosRequestConfig } from 'axios';

import { OAuthToken, SignIn, User } from './authModel';

// A REACT_APP_API_URL value can be injected by webpack during build.
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
const USER_KEY = 'authx.user';
const OAUTH_TOKEN_KEY = 'authx.oauthToken';

const config: AxiosRequestConfig = {
  responseType: 'json',
  headers: {
    'Content-Type': 'application/json',
  },
  transformRequest: (data, headers) => {
    console.log('request', data, headers);
    return JSON.stringify(data);
  },
  transformResponse: (data, headers) => {
    console.log('response', data, headers);
    return data;
  },
};

function getUserFromStorage(): User | null {
  const val = localStorage.getItem(USER_KEY);
  return val ? JSON.parse(val) : null;
}

function isUserInStorage(): boolean {
  return !!getUserFromStorage();
}

function removeUserFromStorage() {
  localStorage.removeItem(USER_KEY);
}

function setUserToStorage(user: User) {
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

function getOAuthTokenFromStorage(): OAuthToken | null {
  const val = localStorage.getItem(OAUTH_TOKEN_KEY);
  return val ? JSON.parse(val) : null;
}

function isOAuthTokenInStorage(): boolean {
  return !!getUserFromStorage();
}

function removeOAuthTokenFromStorage() {
  localStorage.removeItem(OAUTH_TOKEN_KEY);
}

function setOAuthTokenToStorage(authToken: OAuthToken) {
  localStorage.setItem(OAUTH_TOKEN_KEY, JSON.stringify(authToken));
}

async function signIn(username: string, password: string): Promise<OAuthToken | null> {
  return new Promise<OAuthToken | null>((resolve, reject) => {
    const callAPI = async () => {
      try {
        // Call API.
        const url = new URL(API_BASE_URL);
        url.pathname = '/v1/signin';
        const cfg = { ...config };
        const signInData: SignIn = {
          username,
          password,
        };
        const res = await axios.post<OAuthToken>(url.toString(), signInData, cfg);
        const oauthToken = res.data;

        resolve(oauthToken);
      } catch (err) {
        return reject(err);
      }
    };

    callAPI();
  });
}

async function getUserInfo(authToken: OAuthToken): Promise<User | null> {
  return new Promise<User | null>((resolve, reject) => {
    const calledPAI = async () => {
      try {
        const url = new URL(API_BASE_URL);
        url.pathname = '/v1/userinfo';

        if (!authToken.access_token) {
          return reject('null access token');
        }
        const cfg = { ...config };
        cfg.headers.Authorization = `Bearer ${authToken.access_token}`;
        const res = await axios.get<User>(url.toString(), cfg);
        const user = res.data;

        resolve(user);
      } catch (err) {
        return reject(err);
      }
    };

    calledPAI();
  });
}

export {
  getUserInfo,
  getUserFromStorage,
  isUserInStorage,
  removeUserFromStorage,
  getOAuthTokenFromStorage,
  isOAuthTokenInStorage,
  removeOAuthTokenFromStorage,
  setUserToStorage,
  setOAuthTokenToStorage,
  signIn,
};
