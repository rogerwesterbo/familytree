import { lazy, Suspense } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Theme, Flex, Spinner } from '@radix-ui/themes';
import '@radix-ui/themes/styles.css';
import { ThemeProvider, useTheme, AuthProvider } from './contexts';
import { Layout, ProtectedRoute, ScrollToTop } from './components';

// Lazy load pages
const LoginPage = lazy(() => import('./pages/LoginPage'));
const CallbackPage = lazy(() => import('./pages/CallbackPage'));
const DashboardPage = lazy(() => import('./pages/DashboardPage'));
const ProfilePage = lazy(() => import('./pages/ProfilePage'));
const PersonsPage = lazy(() => import('./pages/PersonsPage'));
const RelationshipsPage = lazy(() => import('./pages/RelationshipsPage'));
const RelationshipDetailPage = lazy(() => import('./pages/RelationshipDetailPage'));
const SearchPage = lazy(() => import('./pages/SearchPage').then(m => ({ default: m.SearchPage })));
const ExportPage = lazy(() => import('./pages/ExportPage'));
const GraphViewPage = lazy(() => import('./pages/GraphViewPage'));
const UnauthorizedPage = lazy(() => import('./pages/UnauthorizedPage'));
const ForbiddenPage = lazy(() => import('./pages/ForbiddenPage'));
const NotFoundPage = lazy(() => import('./pages/NotFoundPage'));
const ServerErrorPage = lazy(() => import('./pages/ServerErrorPage'));

const AdminDashboardPage = lazy(() => import('./pages/admin/AdminDashboardPage'));
const CachePage = lazy(() => import('./pages/admin/CachePage'));

// Loading fallback component
function LoadingFallback() {
  return (
    <Flex align="center" justify="center" style={{ minHeight: '100vh', width: '100%' }}>
      <Spinner size="3" />
    </Flex>
  );
}

function AppContent() {
  const { theme } = useTheme();

  return (
    <Theme accentColor="blue" grayColor="slate" radius="medium" appearance={theme}>
      <AuthProvider>
        <Router>
          <ScrollToTop />
          <Suspense fallback={<LoadingFallback />}>
            <Routes>
              {/* Public routes */}
              <Route path="/login" element={<LoginPage />} />
              <Route path="/callback" element={<CallbackPage />} />

              {/* Protected routes with layout */}
              <Route
                path="/"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <DashboardPage />
                    </Layout>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/profile"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <ProfilePage />
                    </Layout>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/persons"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <PersonsPage />
                    </Layout>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/relationships"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <RelationshipsPage />
                    </Layout>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/relationships/:id"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <RelationshipDetailPage />
                    </Layout>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/search"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <SearchPage />
                    </Layout>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/export"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <ExportPage />
                    </Layout>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/graph"
                element={
                  <ProtectedRoute>
                    <Layout noPadding>
                      <GraphViewPage />
                    </Layout>
                  </ProtectedRoute>
                }
              />

              {/* Admin routes */}
              <Route
                path="/admin"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <AdminDashboardPage />
                    </Layout>
                  </ProtectedRoute>
                }
              />
              <Route
                path="/admin/cache"
                element={
                  <ProtectedRoute>
                    <Layout>
                      <CachePage />
                    </Layout>
                  </ProtectedRoute>
                }
              />

              {/* Error routes */}
              <Route path="/unauthorized" element={<UnauthorizedPage />} />
              <Route path="/forbidden" element={<ForbiddenPage />} />
              <Route path="/server-error" element={<ServerErrorPage />} />
              <Route path="*" element={<NotFoundPage />} />
            </Routes>
          </Suspense>
        </Router>
      </AuthProvider>
    </Theme>
  );
}

function App() {
  return (
    <ThemeProvider>
      <AppContent />
    </ThemeProvider>
  );
}

export default App;
