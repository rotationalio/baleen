import { combineReducers } from 'redux';

import Layout from './layout/reducers';
import { vocabularies } from './vocabulary';

export default combineReducers({
    Layout,
    Vocabularies: vocabularies,
});
