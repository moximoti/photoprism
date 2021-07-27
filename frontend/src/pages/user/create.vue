<template>
  <v-form ref="form" dense autocomplete="off" class="p-form-login" accept-charset="UTF-8" @submit.prevent="login">
    <v-card flat tile class="ma-2 application">
      <v-card-actions>
        <v-layout wrap align-top>
          <v-flex xs12 class="pa-2">
            <v-text-field
                v-model="fullname"
                hide-details
                type="text"
                :disabled="loading"
                :label="$gettext('Full Name')"
                browser-autocomplete="off"
                color="secondary-dark"
                class="input-name"
                placeholder="optional"
            ></v-text-field>
          </v-flex>
          <v-flex xs12 class="pa-2">
            <v-text-field
                v-model="username"
                required hide-details
                type="text"
                :disabled="loading"
                :label="$gettext('Username')"
                browser-autocomplete="off"
                color="secondary-dark"
                class="input-name"
                placeholder="username"
            ></v-text-field>
          </v-flex>
          <v-flex xs12 class="pa-2">
            <v-text-field
                v-model="email"
                hide-details
                type="text"
                :disabled="loading"
                :label="$gettext('E-Mail')"
                browser-autocomplete="off"
                color="secondary-dark"
                class="input-name"
                placeholder="me@example.com"
                @keyup.enter.native="register"
            ></v-text-field>
          </v-flex>
          <v-flex xs12 class="pa-2" v-if="!linkUser">
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
            ></v-text-field>
          </v-flex>
          <v-flex xs12 class="pa-2" v-if="!linkUser">
            <v-text-field
                v-model="passwordConfirm"
                required hide-details
                :type="showPassword ? 'text' : 'password'"
                :disabled="loading"
                :label="$gettext('Confirm Password')"
                browser-autocomplete="off"
                color="secondary-dark"
                placeholder="••••••••"
                class="input-password"
                :append-icon="showPassword ? 'visibility' : 'visibility_off'"
                @click:append="showPassword = !showPassword"
                @keyup.enter.native="register"
            ></v-text-field>
          </v-flex>
          <v-flex xs12 class="px-2 py-3">
            <v-btn color="primary-button"
                   class="white--text ml-0 action-confirm"
                   depressed
                   :disabled="!username || !email"
                   @click.stop="register">
              <translate>Complete Registration</translate>
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
  name: 'UserCreate',
  data() {
    const linkUser = window.localStorage.getItem('link_user');
    const userData = JSON.parse(linkUser);

    return {
      linkUser: !!linkUser,
      email: userData?.Email,
      fullname: userData?.Name,
      username: userData?.NickName,
      password: "",
      showPassword: false,
      loading: false,
      nextUrl: this.$route.params.nextUrl ? this.$route.params.nextUrl : "/",
      rtl: this.$rtl,
    };
  },
  methods: {
    register() {
      if (this.linkUser) {
        this.$session.register(this.username, null, null, this.fullname, this.email, this.getIdToken())
          .then(() => {
            this.loading = false;
            window.localStorage.removeItem('link_user');
            this.$router.push(this.nextUrl);
          });
      }
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
