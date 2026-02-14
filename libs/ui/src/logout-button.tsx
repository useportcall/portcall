import { logout } from "./api/logout";

export function LogoutButton() {
  return (
    <form action={logout}>
      <button type="submit">Logout</button>
    </form>
  );
}
