/**
 * This is some documentation
 */
// This is a single line comment
@json
export function foo(player: Player): void {}

/**
 * This is a class that represents a Three Dimensional Vector
 */
@json
class Vec3 {
  x: f32 = 0.0;
  /**
   * This represents the x axis
   */
  y: f32 = 0.0;
  z: f32 = 0.0;
}

/**
 * This class represents a player in a fictitious game
 * @field firstName - ""
 */
@json
class Player {
  @alias("first name")
  firstName!: string;
  lastName!: string;
  /** 
   * This is some docs describing lastActive
  */
  lastActive!: i32[];
  @omitif("this.age < 18")
  age!: i32;
  @omitnull()
  pos!: Vec3 | null;
  isVerified!: boolean;
}