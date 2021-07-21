<template>
  <div class="p-page p-page-login">
    <v-toolbar flat color="secondary" dense class="mb-0" :height="42">
      <v-toolbar-title class="subheading">
        {{ siteDescription }}
      </v-toolbar-title>
    </v-toolbar>
    <v-tabs
        v-model="active"
        flat
        grow
        touchless
        color="secondary-light"
        slider-color="secondary-dark"
        :height="$vuetify.breakpoint.smAndDown ? 48 : 64"
    >
<!--      <v-tab v-for="(item, index) in tabs" :id="'tab-' + item.name" :key="index" :class="item.class" ripple >-->
      <v-tab v-for="(item, index) in tabs" :id="'tab-' + item.name" :key="index" :class="item.class" ripple
             @click="changePath(item.path)">
        {{ item.name }}
      </v-tab>
      <v-tabs-items touchless>
        <v-tab-item v-for="(item, index) in tabs" :key="index" lazy>
          <component :is="item.component"></component>
<!--          <div>{{ item.name }} TEXT TEXT</div>-->
        </v-tab-item>
      </v-tabs-items>
    </v-tabs>

    <p-about-footer></p-about-footer>
  </div>
</template>

<script>
import Login from "pages/user/login.vue";
import Create from "pages/user/create.vue";
function initTabs(flag, tabs) {
  let i = 0;
  while(i < tabs.length) {
    if(!tabs[i][flag]) {
      tabs.splice(i,1);
    } else {
      i++;
    }
  }
}

export default {
  name: 'Register',
  props: {
    tab: String,
  },
  data() {
    const c = this.$config.values;
    const isDemo = this.$config.get("demo");
    const isPublic = this.$config.get("public");
    const tabs = [
      {
        'name': 'register new user',
        'path': '/register',
        'component': Create,
        'public': true,
        'admin': true,
        'demo': true,
        'props': {}
      },
      {
        'name': 'link to existing user',
        'path': '/register',
        'component': Login,
        'public': true,
        'admin': true,
        'demo': true,
        'props': {}
      },
      // {
      //   'name': 'login',
      //   'path': '/login',
      //   'component': Login,
      //   'public': true,
      //   'admin': true,
      //   'demo': true,
      // },
    ];

    if(isDemo) {
      initTabs("demo", tabs);
    } else if(isPublic) {
      initTabs("public", tabs);
    }

    let active = 0;
    if (typeof this.tab === 'string' && this.tab !== '') {
      active = tabs.findIndex((t) => t.name === this.tab);
    }

    return {
      siteDescription: c.siteDescription ? c.siteDescription : c.siteCaption,
      tabs: tabs,
      active: active,
      rtl: this.$rtl,
    };
  },
  methods: {
    changePath: function (path) {
      if (this.$route.path !== path) {
        this.$router.replace(path);
      }
    }
  },
};
</script>
