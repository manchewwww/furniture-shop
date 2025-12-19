import React, { createContext, useContext, useMemo, useState } from "react";

type Lang = "en" | "bg";

type I18nCtx = {
  lang: Lang;
  setLang: (l: Lang) => void;
  t: (k: string) => string;
};

const dict: Record<Lang, Record<string, string>> = {
  en: {
    "nav.home": "Home",
    "nav.catalog": "Catalog",
    "nav.cart": "Cart",
    "nav.orders": "My Orders",
    "nav.admin": "Admin",
    "nav.login": "Login",
    "nav.register": "Register",
    "nav.logout": "Logout",
    "footer.copyright": "Furniture Shop © 2025",

    "home.title": "Welcome to the Furniture Shop",
    "home.subtitle": "Browse departments and discover furniture that fits.",

    "catalog.title": "Catalog",
    "catalog.select.department": "Select department",
    "catalog.select.category": "Select category",
    "catalog.view": "View",

    "cart.title": "Cart",

    "checkout.title": "Checkout",
    "checkout.form.title": "Create Order and Choose Payment",
    "checkout.name": "Name",
    "checkout.email": "Email",
    "checkout.phone": "Phone",
    "checkout.address": "Address",
    "checkout.payment_method": "Payment Method",
    "checkout.place_order": "Place Order",
    "checkout.success": "Order created successfully",
    "checkout.error": "Failed to create order",
    "checkout.card.title": "Pay by Card",
    "checkout.cardholder_name": "Cardholder Name",
    "checkout.card_number": "Card Number",
    "checkout.exp_month": "Month",
    "checkout.exp_year": "Year",
    "checkout.pay": "Pay",
    "checkout.pay.success": "Payment successful",
    "checkout.pay.error": "Payment declined",

    "login.title": "Login",
    "login.email": "Email",
    "login.password": "Password",
    "login.submit": "Login",
    "login.success": "Logged in successfully",
    "login.error": "Invalid email or password",
    "login.register_cta": "Don't have an account?",

    "register.title": "Register",
    "register.name": "Name",
    "register.email": "Email",
    "register.password": "Password",
    "register.address": "Address",
    "register.phone": "Phone",
    "register.submit": "Register",
    "register.success": "Registered successfully",
    "register.error": "Registration failed",

    "orders.title": "My Orders",
    "orders.error": "Please log in to view your orders",
    "orders.col.id": "ID",
    "orders.col.status": "Status",
    "orders.col.payment_status": "Payment Status",
    "orders.col.total": "Total",
    "orders.col.eta_days": "ETA (days)",

    "product.base_price": "Base price",
    "product.base_prod_time": "Base production time (days)",
    "product.options": "Options",
    "product.select_options": "Select options",
    "product.quantity": "Quantity",
    "product.add_to_cart": "Add to cart",
    "product.added": "Added to cart",
    "product.recommended": "Recommended products",

    "admin.departments": "Departments",
    "admin.create_department": "Create Department",
    "admin.department_name": "Name",
    "admin.department_description": "Description",
    "admin.products": "Products",
    "admin.create_product": "Add Product",
    "admin.edit_product": "Edit Product",
    "admin.product_name": "Name",
    "admin.product_description": "Description",
    "admin.product_price": "Price",
    "admin.product_image": "Image (URL)",
    "admin.category": "Category",
    "admin.department": "Department",
    "admin.product_production_days": "Estimated delivery (days)",
    "admin.upload_image": "Upload Image",
    "admin.upload_success": "Uploaded",
    "admin.upload_fail": "Upload failed",
    "admin.actions": "Actions",
    "admin.orders": "Orders",
    "admin.orders_title": "Orders Management",
  },
  bg: {
    "nav.home": "Начало",
    "nav.catalog": "Каталог",
    "nav.cart": "Количка",
    "nav.orders": "Моите поръчки",
    "nav.admin": "Админ",
    "nav.login": "Вход",
    "nav.register": "Регистрация",
    "nav.logout": "Изход",
    "footer.copyright": "Furniture Shop © 2025",
    "home.title": "Добре дошли в магазина за мебели",
    "home.subtitle": "Разгледайте отделите и открийте подходящите мебели.",
    "catalog.title": "Каталог",
    "catalog.select.department": "Изберете отдел",
    "catalog.select.category": "Изберете категория",
    "catalog.view": "Преглед",
    "cart.title": "Количка",
    "checkout.title": "Поръчка",
    "checkout.form.title": "Създаване на поръчка и избор на плащане",
    "checkout.name": "Име",
    "checkout.email": "Имейл",
    "checkout.phone": "Телефон",
    "checkout.address": "Адрес",
    "checkout.payment_method": "Метод на плащане",
    "checkout.place_order": "Направи поръчка",
    "checkout.success": "Поръчката е създадена",
    "checkout.error": "Грешка при създаване на поръчка",
    "checkout.card.title": "Плащане с карта",
    "checkout.cardholder_name": "Име на картодържател",
    "checkout.card_number": "Номер на карта",
    "checkout.exp_month": "Месец",
    "checkout.exp_year": "Година",
    "checkout.pay": "Плати",
    "checkout.pay.success": "Плащането е успешно",
    "checkout.pay.error": "Плащането е отказано",
    "login.title": "Вход",
    "login.email": "Имейл",
    "login.password": "Парола",
    "login.submit": "Вход",
    "login.success": "Успешен вход",
    "login.error": "Грешен имейл или парола",
    "login.register_cta": "Нямате акаунт?",
    "register.title": "Регистрация",
    "register.name": "Име",
    "register.email": "Имейл",
    "register.password": "Парола",
    "register.address": "Адрес",
    "register.phone": "Телефон",
    "register.submit": "Регистрация",
    "register.success": "Успешна регистрация",
    "register.error": "Грешка при регистрация",
    "orders.title": "Моите поръчки",
    "orders.error": "Моля, влезте в профила си",
    "orders.col.id": "№",
    "orders.col.status": "Статус",
    "orders.col.payment_status": "Плащане",
    "orders.col.total": "Общо",
    "orders.col.eta_days": "Срок (дни)",
    "product.base_price": "Базова цена",
    "product.base_prod_time": "Базово време за изработка (дни)",
    "product.options": "Опции",
    "product.select_options": "Изберете опции",
    "product.quantity": "Количество",
    "product.add_to_cart": "Добави в количката",
    "product.added": "Добавено в количката",
    "product.recommended": "Подобни продукти",
    "admin.departments": "Отдели",
    "admin.create_department": "Създай отдел",
    "admin.department_name": "Име",
    "admin.department_description": "Описание",
    "admin.products": "Продукти",
    "admin.create_product": "Добави продукт",
    "admin.edit_product": "Редакция на продукт",
    "admin.product_name": "Име",
    "admin.product_description": "Описание",
    "admin.product_price": "Цена",
    "admin.product_image": "Изображение (URL)",
    "admin.category": "Категория",
    "admin.upload_image": "Качи изображение",
    "admin.upload_success": "Качено",
    "admin.upload_fail": "Качването неуспешно",
    "admin.actions": "Действия",
    "admin.orders": "Поръчки",
    "admin.orders_title": "Управление на поръчки",
  },
};

const I18nContext = createContext<I18nCtx | undefined>(undefined);

export const I18nProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [lang, setLangState] = useState<Lang>(
    () => (localStorage.getItem("lang") as Lang) || "en"
  );
  const setLang = (l: Lang) => {
    setLangState(l);
    localStorage.setItem("lang", l);
  };
  const t = (k: string) => dict[lang][k] ?? k;
  const value = useMemo(() => ({ lang, setLang, t }), [lang]);
  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
};

export const useI18n = () => {
  const ctx = useContext(I18nContext);
  if (!ctx) throw new Error("I18nProvider is missing");
  return ctx;
};
