<template>
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
<!--          <v-flex xs12 class="pa-2">-->
<!--            <v-text-field-->
<!--                v-model="fullname"-->
<!--                hide-details-->
<!--                type="text"-->
<!--                :disabled="loading"-->
<!--                :label="$gettext('Full Name')"-->
<!--                browser-autocomplete="off"-->
<!--                color="secondary-dark"-->
<!--                class="input-name"-->
<!--                placeholder="your name"-->
<!--            ></v-text-field>-->
<!--          </v-flex>-->
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
              <translate>Link User</translate>
              <v-icon :right="!rtl" :left="rtl" dark>login</v-icon>
            </v-btn>
            <v-btn color="primary-button"
                   class="white--text ml-0 action-confirm"
                   depressed
                   :disabled="loading"
                   @click.stop="loginExternal">
              <translate>Forgot Password?</translate>
              <v-icon :right="!rtl" :left="rtl" dark>login</v-icon>
            </v-btn>
          </v-flex>
        </v-layout>
      </v-card-actions>
    </v-card>
  </v-form>
</template>

<script>

export default {
  name: 'UserLink',
  data() {
    const linkUser = window.localStorage.getItem('link_user');
    const userData = linkUser ? JSON.parse(linkUser) : null;

    return {
      linkUser: !!userData,
      username: "",
      password: "",
      showPassword: false,
      loading: false,
      nextUrl: "/",
      rtl: this.$rtl,
    };
  },
  methods: {
    login() {
      if (!this.username || !this.password) {
        return;
      }

      this.loading = true;
      this.$session.login(this.username, this.password, null, this.getIdToken()).then(
        () => {
          this.loading = false;
          window.localStorage.removeItem('link_user');
          this.$router.push(this.nextUrl);
        }
      ).catch(() => this.loading = false);
    },
    getIdToken() {
      const linkUser = window.localStorage.getItem('link_user');
      if (linkUser === null) {
        return null;
      }
      const token = JSON.parse(linkUser).IdToken;
      if (!token) {
        return null;
      }
      return token;
    }
  },
};
</script>
