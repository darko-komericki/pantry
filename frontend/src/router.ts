import type { QueryClient } from '@tanstack/react-query'
import { createRouter } from '@tanstack/react-router'
import { indexRoute } from './routes/index'
import { rootRoute } from './routes/root'

const routeTree = rootRoute.addChildren([indexRoute])

// The QueryClient is passed as router context so route loaders can prefetch
// with queryClient.ensureQueryData() before a page renders.
export function createAppRouter(queryClient: QueryClient) {
  return createRouter({ routeTree, context: { queryClient } })
}

declare module '@tanstack/react-router' {
  interface Register {
    router: ReturnType<typeof createAppRouter>
  }
}
