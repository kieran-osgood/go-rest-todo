import * as React from 'react';

import styles from '../styles/Home.module.scss';
import Layout from '../components/layout';
import Form from '../components/form';
import SearchCard from '../components/search-card';

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
        {/* <SearchCard /> */}
      </main>
    </Layout>
  );
}

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
