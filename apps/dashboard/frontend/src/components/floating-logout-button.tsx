// import { useAuth } from "@clerk/clerk-react";
import { useAuth } from "@/lib/keycloak/auth";
import { Button } from "./ui/button";

export default function FloatingLogoutButton() {
  const { logout } = useAuth();

  return (
    <Button
      variant="ghost"
      className="absolute bottom-10 cursor-pointer hover:text-cyan-800 text-slate-800"
      onClick={() => {
        logout();
      }}
    >
      Logout
    </Button>
  );
}
