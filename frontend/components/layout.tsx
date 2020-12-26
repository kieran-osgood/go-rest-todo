import Head from 'next/head';
import React from 'react';

import styles from '../styles/Home.module.scss';

interface Props {
  children: React.ReactNode
}

const Layout = ({ children }: Props) => (
  <div className={styles.container}>
    <Head>
      <title>Fuzzy Search</title>
      <link rel="icon" href="/favicon.ico" />
    </Head>
    {children}
  </div>
);

export default Layout;
