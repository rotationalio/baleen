import React from 'react';
import PageTitle from 'components/PageTitle';
import Datatable from './Datatable';
import { useDispatch, useSelector } from 'react-redux';
import useSafeDispatch from 'hooks/useSafeDispatch';
import { getAllVocabularies } from 'redux/vocabulary';
import { fetchVocabulariesLoadingState, fetchVocabulariesState } from 'redux/vocabulary/selectors';
import { Vocabulary as VocabularyTypes } from 'types';

const Vocabulary: React.FC = () => {
    const dispatch = useDispatch();
    const safeDispatch = useSafeDispatch(dispatch);
    const data: VocabularyTypes = useSelector(fetchVocabulariesState);
    const isLoading = useSelector(fetchVocabulariesLoadingState);

    React.useEffect(() => {
        safeDispatch(getAllVocabularies());
    }, [safeDispatch]);

    return (
        <div>
            <PageTitle breadCrumbItems={[{ label: 'Vocabulary', path: '/vocabulary' }]} title={'Vocabulary'} />
            {Object.entries(data || []).map(([k, v]) => (
                <Datatable data={v} isLoading={isLoading} title={k} key={k} />
            ))}
        </div>
    );
};

export default Vocabulary;
