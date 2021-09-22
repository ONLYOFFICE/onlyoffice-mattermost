import en from 'i18n/en.json';
import ru from 'i18n/ru.json';
import de from 'i18n/de.json';
import es from 'i18n/es.json';
import it from 'i18n/it.json';
import fr from 'i18n/fr.json';

export function getTranslations(locale?: string) {
    //TODO: Replace with FormattedMessage (at the moment there is a bug with IntlProvider)
    if (locale) {
        window.localStorage.setItem('onlyoffice_locale', locale);
    }
    const currentLocale = locale || window.localStorage.getItem('onlyoffice_locale') || 'en';
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
