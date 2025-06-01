import React from 'react';
import Layout from '../components/Layout';
import AddNovelForm from '../components/AddNovelForm';

const AddNovel = () => {
  return (
    <Layout>
      <div className="container mx-auto py-8 px-4">
        <h1 className="text-3xl font-bold mb-8 text-center text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">Add New Novel</h1>
        <div className="glass-card rounded-lg p-6 max-w-3xl mx-auto">
          <AddNovelForm />
        </div>
      </div>
    </Layout>
  );
};

export default AddNovel;
