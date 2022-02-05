export type Vocabulary = {
    total_documents: number;
    total_unique_words: number;
    total_words: number;
    words_per_document: number;
    most_common_words: {
        [prop: string]: Record<string, any>;
    };
};
