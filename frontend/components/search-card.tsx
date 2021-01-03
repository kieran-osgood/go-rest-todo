import * as React from 'react';
import { useQuery, UseQueryResult } from 'react-query';
import axios, { AxiosError } from 'axios';
import { useDebounce } from 'use-debounce';

import styles from '../styles/Home.module.scss';

const getSearchTodos = async (name: string) => {
  const { data } = await axios.get(
    `/api/todos/search?search=${name}`,
  );
  return data;
};

type DateTime = {Time: string, Valid: boolean};
export type Todo = {
  creationTimestamp: DateTime;
  updateTimestamp: DateTime;
  id: string;
  isDone: boolean;
  text: string;
}

const useSearchTodos = (searchString: string): UseQueryResult<{data: Todo[]}, AxiosError> => useQuery(['search', searchString], () => getSearchTodos(searchString));

export default function SearchCard() {
  const [search, setSearch] = React.useState('');
  const [debouncedSearchTerm] = useDebounce(search, 500);

  const { data: todoList } = useSearchTodos(debouncedSearchTerm);
  return (
    <>
      <div className={styles.card}>
        <p className={styles.description}>
          Fuzzy search/filter:
        </p>

        Search:&nbsp;
        <input type="text" value={search} onChange={(e) => setSearch(e.currentTarget.value)} />
        <ul className={styles.list}>
          {todoList?.data?.map((todo) => (
            <div className={styles.listItem}>
              <input type="checkbox" className={styles.checkbox} />
              <li key={todo.id} className={styles.listText}>{todo.text}</li>
            </div>
          ))}
        </ul>
      </div>
    </>
  );
}
