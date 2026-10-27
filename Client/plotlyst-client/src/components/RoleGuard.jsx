import { Navigate, Outlet, useLocation } from 'react-router-dom';
import useAuth from '../hooks/useAuth';
import Spinner from './spinner/Spinner';

export default function RoleGuard({ roles }) {
  const { auth, loading } = useAuth();
  const location = useLocation();
  if (loading) return <Spinner />;
  if (!auth) return <Navigate to="/login" state={{ from: location }} replace />;
  return roles.includes(auth.role) ? <Outlet /> : <Navigate to="/library" replace />;
}
