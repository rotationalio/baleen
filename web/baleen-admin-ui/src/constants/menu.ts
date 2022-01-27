export interface MenuItemTypes {
    key: string;
    label: string;
    isTitle?: boolean;
    icon?: string;
    url?: string;
    badge?: {
        variant: string;
        text: string;
    };
    parentKey?: string;
    target?: string;
    children?: MenuItemTypes[];
}

const MENU_ITEMS: MenuItemTypes[] = [
    {
        key: 'dashboards',
        label: 'Dashboard',
        isTitle: false,
        icon: 'home',
        url: '/dashboard',
    },
    { key: 'apps', label: 'Apps', isTitle: true },
    {
        key: 'apps-vocab',
        label: 'Vocabulary',
        isTitle: false,
        icon: 'book',
        url: '/vocabulary',
    },
    {
        key: 'apps-topic',
        label: 'Topic',
        isTitle: false,
        icon: 'message-square',
        url: '/topics',
    },
];

export { MENU_ITEMS };
