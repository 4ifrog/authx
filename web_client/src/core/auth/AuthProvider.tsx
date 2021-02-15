import { createContext, ReactNode, useContext } from 'react';

import { User } from './authModel';
import * as authService from './authService';

interface Auth {
  getUser: () => User | null | undefined;
  isSignedIn: () => boolean;
  signIn: (username: string, password: string) => void;
  signOut: () => void;
}

interface AuthProviderProps {
  children?: ReactNode;
}

const defaultAuth = {
  getUser: () => null,
  isSignedIn: () => false,
  signIn: () => {},
  signOut: () => {},
};

const AuthContext = createContext<Auth>(defaultAuth);
const useAuth = () => useContext(AuthContext);

function AuthProvider({ children, ...rest }: AuthProviderProps) {
  let getUser = () => {
    return authService.getUserFromStorage();
  };
  const isSignedIn = () => {
    return !!authService.getOAuthTokenFromStorage();
  };
  const signOut = () => {
    authService.removeUserFromStorage();
    authService.removeOAuthTokenFromStorage();
  };
  const signIn = async (username: string, password: string): Promise<User | null> => {
    return new Promise<User | null>(async (resolve, reject) => {
      try {
        // Get access and refresh tokens (encapsulated in an auth token object).
        const oauthToken = await authService.signIn(username, password);
        if (!oauthToken) {
          return reject(new Error('failed authentication'));
        }
        authService.setOAuthTokenToStorage(oauthToken);

        // Get user with an access token.
        const user = await authService.getUserInfo(oauthToken);
        if (!user) {
          return reject(new Error('failed to fetch user'));
        }

        authService.setUserToStorage(user);
        resolve(user);
      } catch (err) {
        return reject(err);
      }
    });
  };

  const auth: Auth = {
    getUser,
    isSignedIn,
    signOut,
    signIn,
  };

  return (
    <AuthContext.Provider value={auth} {...rest}>
      {children}
    </AuthContext.Provider>
  );
}

export type { Auth, AuthProviderProps };
export { AuthContext, AuthProvider, useAuth };
