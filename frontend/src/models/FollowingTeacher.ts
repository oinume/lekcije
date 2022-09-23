import type {FollowingTeacher as FollowingTeacherGraphQL} from '../graphql/generated';

export type FollowingTeacher = Omit<FollowingTeacherGraphQL, 'id' | 'createdAt' | '__typename'>;
