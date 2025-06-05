import React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-react';

const ITEMS_PER_PAGE_OPTIONS = [10, 20, 50, 100];

interface PaginationControlsProps {
  pagination: {
    novels: any[];
    total_count: number;
    current_page: number;
    total_pages: number;
  };
  loading: boolean;
  itemsPerPage: number;
  pageInput: string;
  onPageChange: (page: number) => void;
  onPageInputChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onPageInputBlur: () => void;
  onPageInputKeyDown: (e: React.KeyboardEvent<HTMLInputElement>) => void;
  onItemsPerPageChange: (itemsPerPage: number) => void;
}

const PaginationControls: React.FC<PaginationControlsProps> = ({
  pagination,
  loading,
  itemsPerPage,
  pageInput,
  onPageChange,
  onPageInputChange,
  onPageInputBlur,
  onPageInputKeyDown,
  onItemsPerPageChange,
}) => {
  if (pagination.total_pages <= 1) return null;

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between gap-4 mt-6">
      <div className="flex items-center space-x-2">
        <span className="text-sm text-muted-foreground">
          Showing {pagination.novels.length} of {pagination.total_count} novels
        </span>
      </div>
      
      <div className="flex items-center space-x-2">
        <Select
          value={String(itemsPerPage)}
          onValueChange={(value) => onItemsPerPageChange(Number(value))}
        >
          <SelectTrigger className="w-[100px] h-9">
            <SelectValue placeholder="Items per page" />
          </SelectTrigger>
          <SelectContent>
            {ITEMS_PER_PAGE_OPTIONS.map((option) => (
              <SelectItem key={option} value={String(option)}>
                {option} / page
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <div className="flex items-center space-x-1">
          <Button
            variant="outline"
            className="h-9 w-9 p-0"
            onClick={() => onPageChange(1)}
            disabled={pagination.current_page === 1 || loading}
          >
            <ChevronsLeft className="h-4 w-4" />
            <span className="sr-only">First page</span>
          </Button>
          <Button
            variant="outline"
            className="h-9 w-9 p-0"
            onClick={() => onPageChange(pagination.current_page - 1)}
            disabled={pagination.current_page === 1 || loading}
          >
            <ChevronLeft className="h-4 w-4" />
            <span className="sr-only">Previous page</span>
          </Button>
          <div className="flex items-center space-x-1">
            <span className="text-sm px-2">Page</span>
            <Input
              className="w-16 h-9 text-center"
              value={pageInput}
              onChange={onPageInputChange}
              onBlur={onPageInputBlur}
              onKeyDown={onPageInputKeyDown}
              disabled={loading}
            />
            <span className="text-sm px-2">of {pagination.total_pages}</span>
          </div>
          <Button
            variant="outline"
            className="h-9 w-9 p-0"
            onClick={() => onPageChange(pagination.current_page + 1)}
            disabled={pagination.current_page === pagination.total_pages || loading}
          >
            <ChevronRight className="h-4 w-4" />
            <span className="sr-only">Next page</span>
          </Button>
          <Button
            variant="outline"
            className="h-9 w-9 p-0"
            onClick={() => onPageChange(pagination.total_pages)}
            disabled={pagination.current_page === pagination.total_pages || loading}
          >
            <ChevronsRight className="h-4 w-4" />
            <span className="sr-only">Last page</span>
          </Button>
        </div>
      </div>
    </div>
  );
};

export default PaginationControls;
