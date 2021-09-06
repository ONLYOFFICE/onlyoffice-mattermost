import en from 'i18n/en.json';
import ru from 'i18n/ru.json';

export function getTranslations(locale?: string) {
    //TODO: Replace with FormattedMessage (at the moment there is a bug with IntlProvider)
    if (locale) {
        window.localStorage.setItem('temp_locale', locale);
    }
    const currentLocale = locale || window.localStorage.getItem('temp_locale');
    switch (currentLocale) {
    case 'de':
        return de;
    case 'en':
        return en;
    case 'es':
        return es;
    case 'fr':
        return fr;
    case 'it':
        return it;
    case 'ru':
        return ru;
    default:
        return en;
    }
}
