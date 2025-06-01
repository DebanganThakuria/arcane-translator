import React from 'react';
import Layout from '../components/Layout';
import { Button } from '@/components/ui/button';
import { Link } from 'react-router-dom';

const NotFound = () => {
  return (
    <Layout>
      <div className="container mx-auto py-16 px-4 text-center">
        <h1 className="text-6xl font-bold mb-4 text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-blue-400">404</h1>
        <p className="text-xl text-muted-foreground mb-6">
          Oops! The page you're looking for can't be found.
        </p>
        <Button asChild className="bg-gradient-to-r from-indigo-600 to-blue-400 hover:from-indigo-700 hover:to-blue-500 border-none">
          <Link to="/">Return to Home</Link>
        </Button>
      </div>
    </Layout>
  );
};

export default NotFound;
