import { RootState } from './../store';
import { createSelector } from 'reselect';

const state = (state: RootState) => state.Vocabularies;

export const fetchVocabulariesState = createSelector(state, (state) => state.data);
export const fetchVocabulariesLoadingState = createSelector(state, (state) => state.isLoading);
