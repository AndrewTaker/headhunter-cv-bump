export interface ErrorResponse {
    error: string;
}

export interface User {
    id: string;
    first_name: string;
    last_name: string;
    middle_name: string;
}

export interface Resume {
    id: string;
    title: string;
    created_at: Date;
    updated_at: Date;
    alternate_url: string;
    is_scheduled: number;
}

export interface ResumeResponseMany {
    resumes: Resume[];
}
