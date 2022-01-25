# Baleen

**An automated ingestion service of RSS feeds to construct a corpus for NLP research.**

[![GoDoc](https://godoc.org/github.com/kansaslabs/baleen?status.svg)](https://godoc.org/github.com/kansaslabs/baleen)
[![Go Report Card](https://goreportcard.com/badge/github.com/kansaslabs/baleen)](https://goreportcard.com/report/github.com/kansaslabs/baleen)

Current overview:

- Golang ingestion system that fetches RSS feeds and stores raw data into LevelDB
- Web-based RSS feed management system that will allow us to easily manage sources
- Focus on fetching full text by following links in the RSS feed
- Feed data quality measurements with language statistics, e.g. words, vocab, etc. rate of corpus growth, number of entities, etc. (we should look at prose for this)
- JSON based logging with limited retention so we don’t fill up our server with logs - tracking of aggregate metrics over time so we know what’s going on and if it's working.
- Produce model based translations for sentences and paragraphs from the source language to target languages; crowdsource feedback by creating an app that allows bilingual users to say if a translation is good or not to establish annotations.
- Annotation quality assessment tools and gamification.
- Periodic checkpoint of data into S3 for archive and analytics and to reduce EC2 expense.
- Estimated cost with 3 yr reserved instance - $64.04 per month (mostly EBS).


## Notes

- [Google I/O 2013 - Advanced Go Concurrency Patterns](https://www.youtube.com/watch?v=QDDwwePbDtw)
