<template>
  <div class="p-page p-page-login">
    <v-toolbar flat color="secondary" dense class="mb-3" :height="42">
      <v-toolbar-title class="subheading">
        {{ siteDescription }}
      </v-toolbar-title>
    </v-toolbar>
    <v-card-title>
      Registration
    </v-card-title>
    <v-form ref="form" dense autocomplete="off" class="p-form-login" accept-charset="UTF-8" @submit.prevent="login">
      <v-card flat tile class="ma-2 application">
        <v-card-actions>
          <v-layout wrap align-top>
            <v-flex xs12 class="pa-2">
              <v-text-field
                  v-model="username"
                  required hide-details
                  type="text"
                  :disabled="loading"
                  :label="$gettext('Name')"
                  browser-autocomplete="off"
                  color="secondary-dark"
                  class="input-name"
                  placeholder="username"
              ></v-text-field>
            </v-flex>
            <v-flex xs12 class="pa-2">
              <v-text-field
                  v-model="password"
                  required hide-details
                  :type="showPassword ? 'text' : 'password'"
                  :disabled="loading"
                  :label="$gettext('Password')"
                  browser-autocomplete="off"
                  color="secondary-dark"
                  placeholder="••••••••"
                  class="input-password"
                  :append-icon="showPassword ? 'visibility' : 'visibility_off'"
                  @click:append="showPassword = !showPassword"
                  @keyup.enter.native="login"
              ></v-text-field>
            </v-flex>
            <v-flex xs12 class="px-2 py-3">
              <v-btn color="primary-button"
                     class="white--text ml-0 action-confirm"
                     depressed
                     :disabled="loading || !password || !username"
                     @click.stop="login">
                <translate>Sign in</translate>
                <v-icon :right="!rtl" :left="rtl" dark>login</v-icon>
              </v-btn>
              <v-btn color="primary-button"
                     class="white--text ml-0 action-confirm"
                     depressed
                     :disabled="loading"
                     @click.stop="loginExternal">
                <translate>Sign in with {{ authProvider }}</translate>
                <v-icon :right="!rtl" :left="rtl" dark>login</v-icon>
              </v-btn>
            </v-flex>
          </v-layout>
        </v-card-actions>
      </v-card>
    </v-form>

    <p-about-footer></p-about-footer>
  </div>
</template>

<script>
export default {
  name: 'Register',
  data() {
    const c = this.$config.values;
    const authProviderFull = () => {
      switch (c.authProvider) {
        case "openid-connect": return "OpenID Connect";
        case "github": return "Github";
        case "google": return "Google";
        default: return "External Provider";
      }
    };
    return {
      loading: false,
      showPassword: false,
      username: "",
      password: "",
      siteDescription: c.siteDescription ? c.siteDescription : c.siteCaption,
      authProvider: authProviderFull(),
      nextUrl: this.$route.params.nextUrl ? this.$route.params.nextUrl : "/",
      rtl: this.$rtl,
    };
  },
  methods: {
    login() {
      if (!this.username || !this.password) {
        return;
      }

      this.loading = true;
      this.$session.login(this.username, this.password).then(
        () => {
          this.loading = false;
          this.$router.push(this.nextUrl);
        }
      ).catch(() => this.loading = false);
    },
    loginExternal() {
      this.loading = true;
      let popup = window.open('api/v1/auth/external', "test");
      const onstorage = window.onstorage;
      window.onstorage = () => {
        const sid = window.localStorage.getItem('session_id');
        const data = window.localStorage.getItem('data');
        const config = window.localStorage.getItem('config');
        const error = window.localStorage.getItem('auth_error');
        if (error === "authentication failed") {
          window.localStorage.removeItem('config');
          window.localStorage.removeItem('auth_error');
          window.onstorage = onstorage;
          return;
        }
        if (sid === null || data === null || config === null) {
          return;
        }
        console.log("sid = ", sid);
        this.$session.setId(sid);
        this.$session.setData(JSON.parse(data));
        this.$session.setConfig(JSON.parse(config));
        //this.$session.sendClientInfo();
        this.loading = false;
        this.$router.push(this.nextUrl);
        popup.close();
        window.localStorage.removeItem('config');
        window.localStorage.removeItem('auth_error');
        window.onstorage = onstorage;
      };
    },
  }
};
</script>
