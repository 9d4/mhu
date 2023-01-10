import {createRouter, createWebHistory} from 'vue-router'
import NewUser from '../pages/NewUser.vue'

const routes = [
    {
        path: "/",
        name: "index",
        component: NewUser,
    },
];

const router = createRouter({
    history: createWebHistory('/'),
    routes,
});

export default router;
