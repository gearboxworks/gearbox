import Vue from 'vue';
import {
  Button,
  Container,
  Header,
  Aside,
  Main,
  Footer,
  Menu,
  Submenu,
  MenuItem,
  MenuItemGroup,
  Dropdown,
  DropdownItem,
  DropdownMenu,
  Tabs,
  TabPane,
  RadioGroup,
  Radio,
  Row,
  Col,
} from 'element-ui';
import lang from 'element-ui/lib/locale/lang/en';
import locale from 'element-ui/lib/locale';

locale.use(lang);

Vue.use(Button);
Vue.use(Container);
Vue.use(Header);
Vue.use(Aside);
Vue.use(Main);
Vue.use(Footer);
Vue.use(Menu);
Vue.use(Submenu);
Vue.use(MenuItem);
Vue.use(MenuItemGroup);
Vue.use(DropdownMenu);
Vue.use(Dropdown);
Vue.use(DropdownItem);
Vue.use(Tabs);
Vue.use(TabPane);
Vue.use(RadioGroup);
Vue.use(Radio);
Vue.use(Row);
Vue.use(Col);
