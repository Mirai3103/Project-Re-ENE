export interface ILive2DModel {
    Version: number
    FileReferences: FileReferences
    Groups: Group[]
    HitAreas: HitArea[]
  }
  
  export interface FileReferences {
    Moc: string
    Textures: string[]
    Physics: string
    Pose: string
    DisplayInfo: string
    Expressions: Expression[]
    Motions: Motions
  }
  
  export interface Expression {
    Name: string
    File: string
  }
  
  export interface Motions {
   [key: string]: MotionItem[]
  }
  
  export interface MotionItem {
    File: string
  }
  
 
  
  export interface Group {
    Target: string
    Name: string
    Ids: string[]
  }
  
  export interface HitArea {
    Id: string
    Name: string
  }
  