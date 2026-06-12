import DashboardView from './DashboardView.vue'

export default {
  title: 'Views/Dashboard',
  component: DashboardView,
}

export const Default = () => ({
  components: { DashboardView },
  template: '<DashboardView />',
})
