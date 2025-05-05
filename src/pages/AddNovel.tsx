
import React from 'react';
import Layout from '../components/Layout';
import AddNovelForm from '../components/AddNovelForm';

const AddNovel = () => {
  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <h1 className="text-3xl font-bold mb-8 text-center">Add New Novel</h1>
        <AddNovelForm />
      </div>
    </Layout>
  );
};

export default AddNovel;
