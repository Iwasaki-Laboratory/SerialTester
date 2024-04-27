import {createApp, DirectiveBinding } from 'vue'
import App from './App.vue'
import './style.css';

var app = createApp(App)


// カスタムディレクティブの定義
app.directive('select-all', {
  mounted(el: HTMLInputElement) {
    el.onfocus = () => {
      el.select();
    };
  }
});

app.mount('#app')