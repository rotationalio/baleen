import apiCore from 'helpers/api';

const api = new apiCore();

export function getVocabularies() {
    return api.get('/vocabularies');
}
