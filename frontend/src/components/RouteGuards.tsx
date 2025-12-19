import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../store/AuthContext";

export const RequireAuth: React.FC<{ children: React.ReactElement }> = ({
  children,
}) => {
  const { isAuthenticated } = useAuth();
  const location = useLocation();
  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }
  return children;
};

export const RequireRole: React.FC<{
  role: string;
  children: React.ReactElement;
}> = ({ role, children }) => {
  const { user } = useAuth();
  const location = useLocation();
  if (!user || user.role !== role) {
    return <Navigate to="/" replace state={{ from: location }} />;
  }
  return children;
};

export const ForbidRole: React.FC<{
  role: string;
  children: React.ReactElement;
}> = ({ role, children }) => {
  const { user } = useAuth();
  if (user && user.role === role) {
    return <Navigate to="/" replace />;
  }
  return children;
};

export const ForbidAuth: React.FC<{ children: React.ReactElement }> = ({
  children,
}) => {
  const { isAuthenticated } = useAuth();
  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }
  return children;
};
