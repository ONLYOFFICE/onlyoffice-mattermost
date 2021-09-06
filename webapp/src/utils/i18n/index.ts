import en from 'i18n/en.json';
import ru from 'i18n/ru.json';

export function getTranslations(locale?: string) {
    //TODO: Replace with FormattedMessage (at the moment there is a bug with IntlProvider)
    if (locale) {
        window.localStorage.setItem('onlyoffice_locale', locale);
    }
    const currentLocale = locale || window.localStorage.getItem('onlyoffice_locale');
    switch (currentLocale) {
    case 'en':
        return en;
    case 'ru':
        return ru;
    default:
        return en;
    }
}
