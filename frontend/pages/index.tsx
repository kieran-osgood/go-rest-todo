import * as React from 'react';
import Head from 'next/head';
import styles from '../styles/Home.module.scss';

export default function Home() {
  return (
    <div className={styles.container}>
      <Head>
        <title>Fuzzy Search</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className={styles.main}>
        <h1 className={styles.title}>
          Fuzzy Search
        </h1>

        <p className={styles.description}>
          Fuzzy search:
        </p>

        <div className={styles.card}>
          Search:
          <input type="text" />
        </div>
      </main>

    </div>
  );
}
