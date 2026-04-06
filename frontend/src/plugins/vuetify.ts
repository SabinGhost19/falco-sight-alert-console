import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi'

export default createVuetify({
  components,
  directives,
  defaults: {
    // Global defaults for Google Cloud look
    VCard: {
      elevation: 0,
      border: true,
      class: 'gc-border rounded-lg'
    },
    VBtn: {
      elevation: 0,
      class: 'text-none font-weight-medium rounded-sm'
    },
    VSheet: {
      elevation: 0,
      border: true,
      class: 'gc-border rounded-lg'
    }
  },
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
          primary: '#1A73E8',    // Google Blue
          secondary: '#E8F0FE',  // Light Blue background for active items
          accent: '#1A73E8',
          error: '#D93025',      // Google Red (Critical)
          warning: '#F29900',    // Google Yellow (Warning)
          info: '#5F6368',       // Google Grey (Text)
          success: '#1E8E3E',    // Google Green (Talon Action)
          background: '#F8F9FA', // Clean grey background
          surface: '#FFFFFF',
        }
      }
    }
  }
})
