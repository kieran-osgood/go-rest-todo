import * as React from 'react';
import { useQuery, UseQueryResult } from 'react-query';
import axios, { AxiosError } from 'axios';
import styles from '../styles/Home.module.scss';
import Layout from '../components/layout';

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
type DateTime = {Time: string, Valid: boolean};
type Todo = {
  CreationTimestamp: DateTime;
  UpdateTimestamp: DateTime;
  ID: string;
  IsDone: boolean;
  Text: string;
}

const useListTodos = (): UseQueryResult<{data: Todo[]}, AxiosError> => useQuery('search', () => getPost());

export default function Home() {
  const [search, setSearch] = React.useState('');
  const { data: todoList } = useListTodos();
  return (
    <Layout>
      <main className={styles.main}>
        <h1 className={styles.title}>
          Fuzzy Search
        </h1>
        <p className={styles.description}>
          Fuzzy search:
        </p>

        <div className={styles.card}>
          Search:&nbsp;
          <input type="text" value={search} onChange={(e) => setSearch(e.currentTarget.value)} />
          <ul>
            {todoList?.data?.map((todo) => (
              <li key={todo.ID}>{todo.Text}</li>
            ))}
          </ul>
        </div>
      </main>
    </Layout>
  );
}
