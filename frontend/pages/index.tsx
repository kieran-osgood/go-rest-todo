import * as React from 'react';
import { useQuery, UseQueryResult } from 'react-query';
import axios, { AxiosError } from 'axios';
import { useDebounce } from 'use-debounce';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';

import styles from '../styles/Home.module.scss';
import Layout from '../components/layout';

const getSearchTodos = async (name: string) => {
  const { data } = await axios.get(
    `/api/todos/search?search=${name}`,
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

const useSearchTodos = (searchString: string): UseQueryResult<{data: Todo[]}, AxiosError> => useQuery(['search', searchString], () => getSearchTodos(searchString));

export default function Home() {
  return (
    <Layout>
      <main className={styles.main}>
        <div>
          <h1 className={styles.title}>
            To-Do List
          </h1>
          <Form />
        </div>
        <SearchCard />
      </main>
    </Layout>
  );
}

const todoForm = z.object({
  todo: z.string().nonempty(),
});
type FormData = z.infer<typeof todoForm>;

const Form = () => {
  const { register, handleSubmit } = useForm({
    resolver: zodResolver(todoForm),
  });

  const submit = (data: FormData) => {
    console.log(data);
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit(submit)}>
      <label htmlFor="todo" className={styles.label}>
        <input type="text" id="todo" name="todo" className={styles.input} ref={register} />
      </label>
      <button type="submit" className={styles.button}>Add</button>
    </form>
  );
};

const SearchCard = () => {
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
              <li key={todo.ID} className={styles.listText}>{todo.Text}</li>
            </div>
          ))}
        </ul>
      </div>
    </>
  );
};

// const getPost = async () => {
//   const { data } = await axios.get(
//     '/api/todos',
//   );
//   return data;
// };

// const useListTodos = (): UseQueryResult<{data: Todo[]}, AxiosError> => {
// return useQuery('search', () => getPost());
// }

// const ListTodos = () => {
//   const { data: todoList } = useListTodos();

//   return (
//     <div>
//       <ul>
//         {todoList?.data?.map((todo) => (
//           <li key={todo.ID}>{todo.Text}</li>
//         ))}
//       </ul>
//     </div>
//   );
// };
