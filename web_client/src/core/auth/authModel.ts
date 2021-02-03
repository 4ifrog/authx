interface AuthToken {
  access_token?: string;
  expiry?: Date;
  refresh_token?: string;
}

interface SignIn {
  username: string;
  password: string;
}

interface User {
  id: string;
  username: string;
  full_name?: string;
  avatar_url?: string;
}

export type { AuthToken, SignIn, User };
