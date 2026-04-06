import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi'

export default createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: { mdi }
  },
  theme: {
    defaultTheme: 'googleCloudTheme',
    themes: {
      googleCloudTheme: {
        dark: false,
        colors: {
          primary: '#4285F4',    // Google Blue
          secondary: '#34A853',  // Google Green
          accent: '#FBBC05',     // Google Yellow
          error: '#EA4335',      // Google Red
          info: '#202124',       // Dark text (Google Cloud header style)
          success: '#34A853',
          warning: '#FBBC05',
          background: '#F8F9FA'  // Clean grey background
        }
      }
    }
  }
})
