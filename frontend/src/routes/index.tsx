import { createRoute } from '@tanstack/react-router'
import { rootRoute } from './root'

export const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: HomePage,
})

function HomePage() {
  return <p>Nothing here yet.</p>
}
