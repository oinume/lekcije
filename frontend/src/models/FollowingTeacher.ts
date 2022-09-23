// export class FollowingTeacher {
// //  [key: string]: any
//
//   constructor(
//     public teacher: Teacher,
//   ) {}
// }

import {FollowingTeacher as FollowingTeacherGraphQL} from '../graphql/generated';

export type FollowingTeacher = Omit<FollowingTeacherGraphQL, "id" | "createdAt" | "__typename">
