import { useEffect, useState } from "react";
import { AlertTriangle, BadgeCheck, Mail, RefreshCcw, Shield, UserRound } from "lucide-react";

import { getApiErrorMessage, getCurrentUser } from "../services/authService";
import { useAuth, type User } from "../hooks/useAuth";

function getRoleLabel(role?: string): string {
  if (!role) {
    return "Customer";
  }

  const labels: Record<string, string> = {
    admin: "Admin",
    customer: "Customer",
    user: "Customer",
  };

  return labels[role] ?? role;
}

function getActiveLabel(isActive?: boolean): string {
  if (isActive === false) {
    return "Inactive";
  }

  return "Active";
}

function ProfileSkeleton() {
  return (
    <section className="profile-page" aria-label="Loading profile">
      <div className="profile-skeleton hero" />
      <div className="profile-skeleton panel" />
    </section>
  );
}

function Profile() {
  const { user: authUser } = useAuth();

  const [profile, setProfile] = useState<User | null>(authUser);
  const [isLoading, setIsLoading] = useState(!authUser);
  const [error, setError] = useState<string | null>(null);

  async function loadProfile() {
    setIsLoading(true);
    setError(null);

    try {
      const result = await getCurrentUser();
      setProfile(result);
    } catch (loadError) {
      setError(
        getApiErrorMessage(
          loadError,
          "Profil belum bisa dimuat. Coba lagi sebentar.",
        ),
      );
    } finally {
      setIsLoading(false);
    }
  }

  useEffect(() => {
    let isActive = true;

    async function loadInitialProfile() {
      setError(null);

      try {
        const result = await getCurrentUser();

        if (isActive) {
          setProfile(result);
        }
      } catch (loadError) {
        if (isActive) {
          setError(
            getApiErrorMessage(
              loadError,
              "Profil belum bisa dimuat. Coba lagi sebentar.",
            ),
          );
        }
      } finally {
        if (isActive) {
          setIsLoading(false);
        }
      }
    }

    loadInitialProfile();

    return () => {
      isActive = false;
    };
  }, []);

  if (isLoading) {
    return <ProfileSkeleton />;
  }

  return (
    <section className="profile-page" aria-labelledby="profile-title">
      <header className="profile-hero">
        <div>
          <span className="products-eyebrow">Customer Profile</span>
          <h1 className="profile-title" id="profile-title">
            Akun kamu.
          </h1>
          <p className="profile-copy">
            View your customer identity, email, role, and account status.
          </p>
        </div>

        <span className="profile-avatar" aria-hidden="true">
          <UserRound className="h-10 w-10" />
        </span>
      </header>

      {error && (
        <div className="profile-notice" role="alert">
          <AlertTriangle className="h-5 w-5" aria-hidden="true" />
          <span>{error}</span>
          <button type="button" onClick={loadProfile}>
            <RefreshCcw className="h-4 w-4" aria-hidden="true" />
            Retry
          </button>
        </div>
      )}

      {!profile ? (
        <div className="profile-empty">
          <div>
            <AlertTriangle className="mx-auto mb-3 h-10 w-10" aria-hidden="true" />
            <h2>Profil tidak tersedia.</h2>
            <p>Data akun belum bisa ditampilkan saat ini.</p>
          </div>
        </div>
      ) : (
        <div className="profile-shell">
          <section className="profile-card" aria-labelledby="profile-info-title">
            <div className="profile-card-heading">
              <UserRound className="h-5 w-5" aria-hidden="true" />
              <h2 id="profile-info-title">Identity</h2>
            </div>

            <div className="profile-field-grid">
              <div className="profile-field">
                <span>Name</span>
                <strong>{profile.name || "Nama belum tersedia"}</strong>
              </div>

              <div className="profile-field">
                <span>Email</span>
                <strong>{profile.email}</strong>
              </div>

              <div className="profile-field">
                <span>Role</span>
                <strong>{getRoleLabel(profile.role)}</strong>
              </div>

              <div className="profile-field">
                <span>Status</span>
                <strong>{getActiveLabel(profile.is_active)}</strong>
              </div>
            </div>
          </section>

          <aside className="profile-card accent">
            <div className="profile-card-heading">
              <BadgeCheck className="h-5 w-5" aria-hidden="true" />
              <h2>Account card</h2>
            </div>

            <div className="profile-badge-list">
              <span>
                <Mail className="h-4 w-4" aria-hidden="true" />
                {profile.email}
              </span>
              <span>
                <Shield className="h-4 w-4" aria-hidden="true" />
                {getRoleLabel(profile.role)}
              </span>
            </div>

            <p className="profile-help-copy">
              Profile editing is not available yet. Contact support if your name
              or email needs to be corrected.
            </p>
          </aside>
        </div>
      )}
    </section>
  );
}

export default Profile;