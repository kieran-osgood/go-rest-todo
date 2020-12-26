import * as React from 'react';
import { useQuery } from 'react-query';
import axios from 'axios';
import styles from '../styles/Home.module.scss';
import Layout from '../components/layout';

const baseUrl = process.env.NEXT_PUBLIC_API_URL;

// const getPostById = async (postId) => {
//   const { data } = await axios.get(
//     `${baseUrl}${postId}`,
//   );
//   return data;
// };

const getPost = async () => {
  const { data } = await axios.get(
    '/api/todos',
  );
  return data;
};

export default function Home() {
  const [search, setSearch] = React.useState('');
  // useQuery(['searchById', search], () => getPostById(search));
  const { data } = useQuery(['search'], () => getPost());
  console.log('data: ', data);

  return (
    <Layout>
      <main className={styles.main}>
        <h1 className={styles.title}>
          Fuzzy Search
        </h1>
        {`${baseUrl}/todos`}
        <p className={styles.description}>
          Fuzzy search:
        </p>

        <div className={styles.card}>
          Search:&nbsp;
          <input type="text" value={search} onChange={(e) => setSearch(e.currentTarget.value)} />
          <ul>
            <li>test</li>
          </ul>
        </div>
      </main>
    </Layout>
  );
}
